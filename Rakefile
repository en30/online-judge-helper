$LOAD_PATH.unshift File.expand_path('../lib', __FILE__)
require 'yaml'
require 'online_judge_helper'
require 'hashie'
require 'open3'
require 'open-uri'
require 'colorize'

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

task install: 'install:all'
namespace :install do
  task all: %i{ruby jquery extension}

  task :ruby do
    if Gem::Version.create(RUBY_VERSION) < Gem::Version.create('2.1')
      abort('Online Judge Helepr needs ruby >= 2.1')
    end
  end

  task :jquery do
    puts 'Downloading jquery for chrome_extension...'.yellow
    path = File.expand_path('chrome_extension/jquery.js')
    File.write(
      path,
      open('http://code.jquery.com/jquery-2.1.3.min.js').read
    ) unless File.exist?(path)
    print "done\n\n".yellow
  end

  task :extension do
    puts '=> Drang and drop the chrome_extension directory to Google Chrome to install the extension'.bold
    system(%q{command osascript -e 'tell application "Google Chrome" to open location "chrome://chrome/extensions"'})
    system('open .')
  end
end
