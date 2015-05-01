require 'sinatra'
require 'sinatra/config_file'
require 'sinatra/reloader'
require 'erb'
require 'pry'
require 'fileutils'

config_file_path = File.expand_path('../config.yml', __FILE__)
config_file config_file_path
also_reload config_file_path

post '/problem' do
  site = File.basename(params['site']) || 'unknown'
  time_limit = params['time_limit'] || 2
  id = File.basename(params['id'])
  editor = settings.editor || ENV['EDITOR']

  problems_path = File.expand_path("../problems/#{site}", __FILE__)
  tests_path = File.expand_path("../tests/#{site}", __FILE__)
  FileUtils.mkdir_p problems_path
  FileUtils.mkdir_p tests_path
  problem_file = "#{problems_path}/#{id}.#{settings.default_lang}"
  test_file = "#{tests_path}/#{id}.rb"

  system "#{editor} #{problem_file} &"

  unless File.exist?(test_file)
    File.open test_file, 'w' do |f|
      f.puts "require_relative '../../test_helper'"
      params['samples'].each_with_index do |(_, sample), i|
        f.puts "test_case title: 'Sample case #{i}', time_limit: #{time_limit}, input: <<INPUT, expected: <<EXPECTED"
        f.puts sample['input']
        f.puts 'INPUT'
        f.puts sample['output']
        f.puts 'EXPECTED'
        f.puts
      end
    end
  end
  200
end
