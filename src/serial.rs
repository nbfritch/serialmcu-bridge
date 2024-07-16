use crate::commands::McuCommand;
use std::{collections::HashMap, path::PathBuf, thread::sleep, time::Duration};

#[derive(Debug)]
pub enum Reading {
    Temperature(f32),
    Humidity(f32),
    Lux(u16),
}

#[derive(Debug)]
pub enum ReadingError {
    InvalidSerialError,
    SerialOpenFailure,
    SerialWriteFailure,
    SerialReadFailure,
    SerialDecodeFailure,
    SerialParseError,
}

#[derive(Debug)]
pub struct CommandError {
    pub error_type: ReadingError,
    pub message: String,
}

pub fn execute_commands(
    port: PathBuf,
    commands: Vec<McuCommand>,
) -> HashMap<McuCommand, Result<Reading, CommandError>> {
    let mut cmd_hash = HashMap::new();
    let port_path_str_result = port.to_str();
    if port_path_str_result.is_none() {
        for cmd in commands.iter() {
            cmd_hash.entry(*cmd).or_insert(Err(CommandError {
                error_type: ReadingError::InvalidSerialError,
                message: "Failed to coerce serial port to filepath".into(),
            }));
        }
        return cmd_hash;
    }
    let port_path = port_path_str_result.unwrap();
    let sp_result = tokio_serial::new(port_path, 115_200)
        .timeout(Duration::from_secs(1))
        .open();
    if let Err(e) = sp_result {
        for cmd in commands.iter() {
            cmd_hash.entry(*cmd).or_insert(Err(CommandError {
                error_type: ReadingError::SerialOpenFailure,
                message: format!("Error while opening serial port: {}", e),
            }));
        }
        return cmd_hash;
    }
    let mut sp = sp_result.unwrap();

    for cmd in commands.iter() {
        println!("{:-}", cmd);
        let cmd_str = format!("{}\r", cmd);
        let write_cmd = cmd_str.as_bytes();
        let write_result = sp.write(write_cmd);
        if let Err(e) = write_result {
            for cmd in commands.iter() {
                cmd_hash.entry(*cmd).or_insert(Err(CommandError {
                    error_type: ReadingError::SerialWriteFailure,
                    message: format!("Error while writing to serial port: {}", e),
                }));
            }
            return cmd_hash;
        }
        let mut read_buffer = vec![0; 128];
        let read_result = sp.read(read_buffer.as_mut_slice());
        if let Err(e) = read_result {
            for cmd in commands.iter() {
                cmd_hash.entry(*cmd).or_insert(Err(CommandError {
                    error_type: ReadingError::SerialReadFailure,
                    message: format!("Error while reading from serial port: {}", e),
                }));
            }
            return cmd_hash;
        }
        let read = read_result.unwrap();

        let s = std::str::from_utf8(&read_buffer[0..read - 1]);
        if let Err(e) = s {
            for cmd in commands.iter() {
                cmd_hash.entry(*cmd).or_insert(Err(CommandError {
                    error_type: ReadingError::SerialDecodeFailure,
                    message: format!("Error while decoding serial port output: {}", e),
                }));
            }
            return cmd_hash;
        }
        let result_string = s.unwrap();
        println!("{:?}", result_string);

        let reading = if *cmd != McuCommand::Lux {
            let float_result = result_string.parse::<f32>();
            if let Err(e) = float_result {
                for cmd in commands.iter() {
                    cmd_hash.entry(*cmd).or_insert(Err(CommandError {
                        error_type: ReadingError::SerialDecodeFailure,
                        message: format!("Error while decoding serial port output: {}", e),
                    }));
                }
                return cmd_hash;
            }
            let float = float_result.unwrap();
            let reading = match cmd {
                McuCommand::Am2320Humidity => Reading::Humidity(float),
                McuCommand::Am2320Temperature => Reading::Temperature(float),
                McuCommand::Ds18x20Temperature => Reading::Temperature(float),
                _ => panic!("Unreachable"),
            };
            reading
        } else {
            let uint_result = result_string.parse::<u16>();
            if let Err(e) = uint_result {
                for cmd in commands.iter() {
                    cmd_hash.entry(*cmd).or_insert(Err(CommandError {
                        error_type: ReadingError::SerialParseError,
                        message: format!("Error while decoding serial port output: {}", e),
                    }));
                }
                return cmd_hash;
            }
            let uint = uint_result.unwrap();
            Reading::Lux(uint)
        };

        cmd_hash.entry(*cmd).or_insert(Ok(reading));
        sleep(Duration::from_millis(500));
    }

    cmd_hash
}
