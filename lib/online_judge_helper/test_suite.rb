require 'erb'
require 'unindent'
module OnlineJudgeHelper
  class TestSuite
    RESULT_TEMPLATE = <<-ERB.unindent
      <%= Time.now %>
      ==============================
      Problem: <%= @title.blue %>
      ==============================
      <% @test_cases.each do |tc| %>
      <%= tc.title.yellow %>: <%= tc.result %>
      <%= tc.message %>

      ------------------------------
      <% end %>
      <%= message %>
    ERB

    def initialize(title:, config:, time_limit: 2, samples:)
      @config = config
      @title = title
      @time_limit = time_limit
      @test_cases = samples.map { |tc| TestCase.new(**tc) }
    end

    def run!
      compile if @config.compile?
      return if compilation_error?

      @test_cases.each do |tc|
        tc.run!(@config.run, @config.problem_file, time_limit: @time_limit)
      end
    end

    def result
      if compilation_error?
        message
      else
        ERB.new(RESULT_TEMPLATE).result(binding)
      end
    end

    def failed?
      compilation_error? or wrong_anser?
    end

    def message
      if compilation_error?
        'Compilation Error'.red + "\n#{@error_message}"
      elsif wrong_anser?
        "Wrong Answer\n#{failed_count}/#{@test_cases.count} failed".red
      else
        'All test passed!!'.green
      end
    end

    private

    def failed_count
      @test_cases.map(&:passed?).count(false)
    end

    def compile
      _, stderr, status = Open3.capture3("#{config.compile} #{config.problem_file}")
      @compilation_error = (status != 0)
      @error_message = stderr
    end

    def compilation_error?
      @compilation_error
    end

    def wrong_anser?
      failed_count != 0
    end
  end
end
