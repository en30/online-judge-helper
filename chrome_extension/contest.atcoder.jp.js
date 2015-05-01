(function($){
  $.appendSolveButton().on('click', function() {
    var params ={
      site: 'atcoder',
      time_limit: $('#task-statement').prev().text().match(/(\d+)sec/)[1],
      id: location.pathname.split('/').pop(),
      samples: []
    };
    $("h3:contains('入力例')").each(function(){
      var $input = $(this).next('pre'),
          $output = $input.closest('.part').next().find('pre');
      params.samples.push({input: $input.text(), output: $output.text()});
    });
    $.solve(params);
  });
})(jQuery);
