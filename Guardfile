
ignore /#/

def run_test(path, lang = '')
  res = `bundle exec ruby #{path} #{lang}`
  case $?.exitstatus
  when 0
    TerminalNotifier::Guard.success 'All tests passed!'
  when 1
    TerminalNotifier::Guard.failed 'Some tests failed'
  when 2
    TerminalNotifier::Guard.failed 'Compilation Error'
  end
  res
end

guard :shell do
  watch(/^problems\/(.+)\.([\w]+)$/) { |m| run_test "tests/#{m[1]}.rb", m[2] }
  watch(/^tests\/.+.rb$/) { |m| run_test(m[0]) }
end
