require 'sinatra/base'
require 'sinatra/config_file'
require 'sinatra/reloader'
require 'pry'
require 'fileutils'
require 'hashie'
require 'open3'

module OnlineJudgeHelper
  class App < Sinatra::Base
    register Sinatra::ConfigFile
    register Sinatra::Reloader

    CONFIG_FILE_PATH = File.expand_path('./config.yml')
    config_file CONFIG_FILE_PATH
    also_reload CONFIG_FILE_PATH

    post '/problem' do
      site = File.basename(params['site']) || 'unknown'
      time_limit = (params['time_limit'] || 2).to_f
      id = File.basename(params['id'])
      editor = settings.editor || ENV['EDITOR']

      problems_path = File.expand_path("./problems/#{site}")
      tests_path = File.expand_path("./tests/#{site}")
      FileUtils.mkdir_p problems_path
      FileUtils.mkdir_p tests_path
      problem_file = "#{problems_path}/#{id}.#{settings.default_language}"
      test_file = "#{tests_path}/#{id}.yml"

      system "#{editor} #{problem_file} &"

      samples = params['samples'].map  do |i, sample|
        { title: "Sample #{i}" }.merge(Hash[sample.map{|k,v| [k.to_sym, v] }])
      end
      File.open(test_file, 'w') do |f|
        f.print(YAML.dump(time_limit: time_limit, samples: samples))
      end
      200
    end
  end
end
