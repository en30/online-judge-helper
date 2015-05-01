require 'open3'
require 'timeout'
require 'colorize'
require 'yaml'

class TestCase
  attr_reader :title, :time_limit, :input, :expected, :actual, :passed
  def initialize(title:, time_limit: 2, input:, expected:)
    @title = title
    @input = input
    @time_limit = time_limit
    @expected = expected
    @passed = false
  end

  def run
    err = '', status = 0
    timeout(time_limit) do
      @actual, err, status = Open3.capture3(config['run'], problem_file, stdin_data: input)
    end
    if status == 0
      @passed = (actual == expected)
      detail = "\nExpected\n#{expected}\n\nbut actual\n\n#{actual}" unless passed
    else
      detail = "Runtime Error: #{status}\n#{actual}\n#{err}"
    end
    passed
  rescue Timeout::Error
    detail = 'Time Limit Exceeded'
    passed
  ensure
    print "#{title.yellow}: "
    if passed
      puts 'PASSED'.green
    else
      puts 'FAILED'.red
      puts detail
    end
    puts '=' * 30
  end
end

@test_cases = []

def lang
  @lang ||= ARGV[0] || YAML.load_file(File.expand_path('../config.yml', __FILE__))['default_lang']
end

def config
  @config ||= YAML.load_file(File.expand_path('../config.yml', __FILE__))['languages'][lang]
end

def test_id
  $0.slice(/tests\/(.*?)\.rb/, 1)
end

def compile?
  config['compile'] and !config['compile'].empty?
end

def problem_file
  File.expand_path("../problems/#{test_id}.#{lang}", __FILE__)
end

def compile
  unless system "#{config['compile']} #{problem_file}"
    puts 'Compilation Failed'.red
    exit 2
  end
end

def test_case(**attrs)
  @test_cases.push TestCase.new(**attrs)
end

at_exit do
  print "\n\n\n"
  puts Time.now
  puts '=' * 30
  puts "Problem: #{test_id.blue}"
  puts '=' * 30
  compile if compile?
  failed_count = @test_cases.map(&:run).count(false)
  if failed_count == 0
    puts 'All test passed!!'.green
  else
    puts "#{failed_count}/#{@test_cases.count} failed".red
    exit 1
  end
end
