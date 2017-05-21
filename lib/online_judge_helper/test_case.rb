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
      @stdout, @stderr = '', ''
      Open3.popen3(command, problem_file) do |i,o,e,w|
        begin
          Timeout.timeout(time_limit) do
            i.write(input)
            i.close
            while !(o.eof? && e.eof?)
              @stdout << o.read(4096).to_s
              @stderr << e.read(4096).to_s
            end
            @status = w.value.exitstatus
          end
        rescue Timeout::Error
          Process.kill('INT', w.pid)
          @result_code = TIME_OUT_ERROR
          return
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
    end
  end
end
