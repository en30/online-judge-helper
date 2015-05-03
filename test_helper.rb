$LOAD_PATH.unshift File.expand_path('../lib', __FILE__)
require 'yaml'
require 'online_judge_helper'
require 'hashie'
require 'open3'

def test_id
  $0.slice(/tests\/(.*?)\.rb/, 1)
end

def config
  return @config if @config
  config = YAML.load_file(File.expand_path('../config.yml', __FILE__))
  config['language'] = ARGV[0] || config['default_lang']
  config['problem_file'] = File.expand_path("../problems/#{test_id}.#{config['language']}", __FILE__)
  config.merge!(config['languages'][config['language']])
  @config = Hashie::Mash.new(config)
end

at_exit do
  test_suite = OnlineJudgeHelper::TestSuite.new(**YAML.load(DATA).merge(config: config, title: test_id))

  test_suite.run!
  puts
  puts test_suite.result

  if test_suite.failed?
    STDERR.puts test_suite.message.uncolorize
    exit 1
  end
end
