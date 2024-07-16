mod args;
mod commands;
mod serial;
mod types;

use args::Args;
use clap::Parser;
use commands::McuCommand;
use serial::execute_commands;
use types::{ReadingType, Sensor};

#[tokio::main]
async fn main() {
    let args = Args::parse();
    let reading_type = if let None = args.reading_type {
        None
    } else {
        let r = ReadingType::try_from(args.reading_type.unwrap());
        match r {
            Ok(a) => Some(a),
            Err(e) => {
                panic!("Error: {}", e);
            }
        }
    };
    let sensor = if let None = args.just_sensor {
        None
    } else {
        let s = Sensor::try_from(args.just_sensor.unwrap());
        match s {
            Ok(a) => Some(a),
            Err(e) => {
                panic!("Error: {}", e);
            }
        }
    };
    let commands = McuCommand::parse(reading_type, sensor);
    match commands {
        Ok(c) => {
            let cmd_result = execute_commands(args.port, c);
            cmd_result.values().for_each(|v| {
                println!("{:?}", v);
            });
        }
        Err(msg) => {
            panic!("{}", msg);
        }
    }
}
