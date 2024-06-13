import serial
import time

read_am2320_hum_cmd = "am2320_humidity"
read_am2320_temp_cmd = "am2320_temp"
read_ds_temp_cmd = "ds18x20_temp"
read_lux_cmd = "lux"

def read_cmd(s, cmd):
    s.flush()
    s.write(f"{cmd}\r".encode())
    time.sleep(1)
    mes = s.read_until().strip()
    return mes.decode()

def main():
    s = serial.Serial(port="/dev/ttyACM0", parity=serial.PARITY_EVEN, stopbits=serial.STOPBITS_ONE, timeout=1)
    print(f"AM2320 Humidity: {read_cmd(s, read_am2320_hum_cmd)}")
    print(f"AM2320 Temp: {read_cmd(s, read_am2320_temp_cmd)}")
    print(f"DS Temp: {read_cmd(s, read_ds_temp_cmd)}")
    print(f"Relative lux: {read_cmd(s, read_lux_cmd)}")

if __name__ == "__main__":
    main()
