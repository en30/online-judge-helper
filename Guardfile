require 'open3'
ignore /#/

def run_test(id, lang = '')
  stdout, stderr, status = Open3.capture3("bundle exec rake test[#{id},#{lang}]")
  if status == 0
    TerminalNotifier::Guard.success '', title: 'All tests passed!'
  else
    title, message = stderr.split("\n", 2)
    TerminalNotifier::Guard.failed message, title: title
  end
  stdout
end

guard :shell do
  watch(/^problems\/(.+)\.([\w]+)$/) { |m| run_test m[1], m[2] }
  watch(/^tests\/(.+).yml$/) { |m| run_test(m[1]) }
end
