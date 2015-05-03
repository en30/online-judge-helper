require 'open3'
require 'timeout'

module OnlineJudgeHelper
  class TestCase
    SUCCESS        = 0
    WRONG_ANSWER   = 1
    RUNTIME_ERROR  = 2
    TIME_OUT_ERROR = 3

    attr_reader :title, :input, :output

    def initialize(title:, input:, output:)
      @title = title
      @input = input
      @output = output
    end

    def message
      case @result_code
      when SUCCESS
      when WRONG_ANSWER
        "\nExpected output: \n#{output}\n\nActual output:\n\n#{@stdout}"
      when RUNTIME_ERROR
        "Runtime Error: #{@status}\n#{@stdout}\n#{@stderr}"
      when TIME_OUT_ERROR
        'Time Limit Exceeded'
      end
    end

    def result
      if passed?
        'PASSED'.green
      else
        'FAILED'.red
      end
    end

    def passed?
      @result_code == SUCCESS
    end

    def run!(command, problem_file, time_limit:)
      # https://redmine.ruby-lang.org/issues/4681
      timeout(time_limit) do
        timeout(time_limit) do
          @stdout, @stderr, @status = Open3.capture3(command, problem_file, stdin_data: input)
        end
      end

      @result_code =
        if @status != 0
          RUNTIME_ERROR
        elsif @stdout == output
          SUCCESS
        else
          WRONG_ANSWER
        end
    rescue Timeout::Error
      @result_code = TIME_OUT_ERROR
    end
  end
end
