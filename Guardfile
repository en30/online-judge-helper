require 'open3'
ignore /#/

def run_test(path, lang = '')
  stdout, stderr, status = Open3.capture3("bundle exec ruby #{path} #{lang}")
  if status == 0
    TerminalNotifier::Guard.success '', title: 'All tests passed!'
  else
    title, message = stderr.split("\n", 2)
    TerminalNotifier::Guard.failed message, title: title
  end
  stdout
end

guard :shell do
  watch(/^problems\/(.+)\.([\w]+)$/) { |m| run_test "tests/#{m[1]}.rb", m[2] }
  watch(/^tests\/.+.rb$/) { |m| run_test(m[0]) }
end
