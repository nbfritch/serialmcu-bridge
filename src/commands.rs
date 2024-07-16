use crate::types::{ReadingType, Sensor};
use std::fmt::Display;

#[derive(PartialEq, Eq, Hash, Copy, Clone, Debug)]
pub enum McuCommand {
    Am2320Humidity,
    Am2320Temperature,
    Ds18x20Temperature,
    Lux,
}

impl Display for McuCommand {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        f.write_str(match self {
            Self::Am2320Humidity => "am2320_humidity",
            Self::Am2320Temperature => "am2320_temp",
            Self::Ds18x20Temperature => "ds18x20_temp",
            Self::Lux => "lux",
        })
    }
}

impl McuCommand {
    pub fn parse(
        reading_type: Option<ReadingType>,
        sensor: Option<Sensor>,
    ) -> Result<Vec<Self>, String> {
        match (reading_type, sensor) {
            (None, None) => Ok(vec![
                Self::Am2320Humidity,
                Self::Am2320Temperature,
                Self::Ds18x20Temperature,
                Self::Lux,
            ]),
            (Some(ReadingType::Temperature), None) => {
                Ok(vec![Self::Am2320Temperature, Self::Ds18x20Temperature])
            }
            (Some(ReadingType::Humidity), None) => Ok(vec![Self::Am2320Humidity]),
            (Some(ReadingType::Lux), None) => Ok(vec![Self::Lux]),
            (None, Some(Sensor::Am2320)) => Ok(vec![Self::Am2320Humidity, Self::Am2320Temperature]),
            (None, Some(Sensor::Ds18x20)) => Ok(vec![Self::Ds18x20Temperature]),
            (None, Some(Sensor::Lux)) => Ok(vec![Self::Lux]),
            _ => Err("invalid query".into()),
        }
    }
}
