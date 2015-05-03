$LOAD_PATH.unshift File.expand_path('../lib', __FILE__)
require 'yaml'
require 'online_judge_helper'
require 'hashie'
require 'open3'

def config
  return @config if @config
  @config = Hashie::Mash.new(YAML.load_file(File.expand_path('config.yml')))
end

desc 'Run a test specified by argument'
task :test, %w{id language} do |_, args|
  test_file = File.expand_path("tests/#{args.id}.yml")
  config.language = args.language.empty? ? config.default_lang : args.language
  config.problem_file = File.expand_path("problems/#{args.id}.#{config.language}")
  config.merge!(config.languages[config.language])

  test_suite = OnlineJudgeHelper::TestSuite.new(**YAML.load_file(test_file).merge(config: config, title: args.id))

  test_suite.run!
  puts
  puts test_suite.result

  if test_suite.failed?
    STDERR.puts test_suite.message.uncolorize
    exit 1
  end
end
