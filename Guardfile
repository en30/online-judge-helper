require 'open3'
ignore /#/
$stdout.sync = true

def run_test(id, lang = '')
  stdout, stderr, status = Open3.capture3("bundle exec rake test[#{id},#{lang}]")
  if status == 0
    TerminalNotifier::Guard.success '', title: 'All tests passed!', timeout: 3
    stdout
  else
    title, message = stderr.split("\n", 2)
    TerminalNotifier::Guard.failed message, title: title, timeout: 3
    stdout + stderr
  end
end

guard :shell do
  watch(/^problems\/(.+)\.([\w]+)$/) { |m| run_test m[1], m[2] }
  watch(/^tests\/(.+).yml$/) { |m| run_test(m[1]) }
end
