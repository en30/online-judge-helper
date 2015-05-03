require 'sinatra/base'
require 'sinatra/config_file'
require 'sinatra/reloader'
require 'pry'
require 'fileutils'

module OnlineJudgeHelper
  class App < Sinatra::Base
    register Sinatra::ConfigFile
    register Sinatra::Reloader

    CONFIG_FILE_PATH = File.expand_path('./config.yml')
    config_file CONFIG_FILE_PATH
    also_reload CONFIG_FILE_PATH

    post '/problem' do
      site = File.basename(params['site']) || 'unknown'
      time_limit = params['time_limit'] || 2
      id = File.basename(params['id'])
      editor = settings.editor || ENV['EDITOR']

      problems_path = File.expand_path("./problems/#{site}")
      tests_path = File.expand_path("./tests/#{site}")
      FileUtils.mkdir_p problems_path
      FileUtils.mkdir_p tests_path
      problem_file = "#{problems_path}/#{id}.#{settings.default_lang}"
      test_file = "#{tests_path}/#{id}.rb"

      system "#{editor} #{problem_file} &"

      unless File.exist?(test_file)
        samples = params['samples'].map  do |i, sample|
          { title: "Sample #{i}" }.merge(Hashie.symbolize_keys sample)
        end
        File.open test_file, 'w' do |f|
          f.puts "require_relative '../../test_helper'"
          f.puts '__END__'
          f.puts YAML.dump(time_limit: time_limit, samples: samples)
        end
      end
      200
    end
  end
end
