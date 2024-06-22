require 'optparse'

options = {}
OptionParser.new do |parser|
  parser.banner = "Usage: serialmcu-bridge.rb [options]"
  parser.on("-n", "--name [NAME]", String)
  parser.on("-p", "--port [PORT]", String)
  parser.on("-r", "--retry [RETRY]", Integer)
  parser.on("-e", "--endpoint [ENDPOINT]", String)
  parser.on("-s", "--send")
  parser.on("-v", "--verbose")
end.parse!(into: options)

AM_TEMP = "am2320_temp"
AM_HUM = "am2320_humidity"
DS_TEMP = "ds18x20_temp"
LUX = "lux"

if options[:port].nil?
  puts "PORT is required"
  return
end

def reading_valid?(reading, type)
  return true if type == :temperature and reading > 0 and reading < 120
  return true if type == :humidity and reading > 0 and reading < 100
  return true if type == :lux and reading > 0 and reading < 65336
  return false
end

def exec_cmd(cmd, port, verbose)
  if verbose
    puts "echo -ne 'r#{cmd}\r' | picocom -qrx 1000 #{port}"
  end
  r = `echo -ne '\r#{cmd}\r' | picocom -qrx 1000 #{port}`
  if verbose
    puts "=> #{r}"
  end
  r.chomp.to_f
end

results = {}
[AM_TEMP, AM_HUM, DS_TEMP, LUX].each do |cmd|
  results[cmd] = exec_cmd(cmd, options[:port], options[:verbose])
end

puts results
