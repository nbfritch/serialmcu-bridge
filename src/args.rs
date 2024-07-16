use std::path::PathBuf;

use clap::Parser;

#[derive(Parser, Debug)]
#[command(version, about, long_about = None)]
pub struct Args {
    #[arg(short, long, value_name = "SERIAL_PORT")]
    pub port: PathBuf,

    #[arg(short, long)]
    pub reading_type: Option<String>,

    #[arg(short, long, value_name = "SENSOR")]
    pub just_sensor: Option<String>,

    #[arg(short, long, default_value_t = false)]
    pub send: bool,
}
