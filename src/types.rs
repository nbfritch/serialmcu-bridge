pub enum Sensor {
    Am2320,
    Ds18x20,
    Lux,
}

impl TryFrom<String> for Sensor {
    type Error = String;
    fn try_from(value: String) -> Result<Self, Self::Error> {
        match value.as_str() {
            "am2320" => Ok(Self::Am2320),
            "ds18x20" => Ok(Self::Ds18x20),
            "lux" => Ok(Self::Lux),
            _ => Err(format!(
                "Invalid sensor {}, expected one of am2320, ds18x20, lux",
                value
            )),
        }
    }
}

pub enum ReadingType {
    Temperature,
    Humidity,
    Lux,
}

impl TryFrom<String> for ReadingType {
    type Error = String;
    fn try_from(value: String) -> Result<Self, Self::Error> {
        match value.as_str() {
            "temperature" => Ok(Self::Temperature),
            "humidity" => Ok(Self::Humidity),
            "lux" => Ok(Self::Lux),
            _ => Err(format!(
                "Invalid reading_type {}, expected one of temperature/temp, humidity/hum, lux",
                value
            )),
        }
    }
}
