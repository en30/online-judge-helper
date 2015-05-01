(function($){
  $.appendSolveButton().on('click', function() {
    var params ={
      site: 'codeforces',
      time_limit: $('.time-limit').text().match(/(\d+) seconds/)[1],
      id: location.pathname.split('/').slice(3).join('_'),
      samples: []
    };
    $(".sample-test > .input").each(function(){
      var $input = $(this).find('pre'),
          $output = $(this).next().find('pre'),
          br2nl = function(html){
            return html.replace(/<br>/g, "\n");
          };
      params.samples.push({input: br2nl($input.html()), output: br2nl($output.html())});
    });
    $.solve(params);
  });
})(jQuery);
