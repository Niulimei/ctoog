 layui.use('element', function() {
     var element = layui.element; //导航的hover效果、二级菜单等功能，需要依赖element模块
     var $iframe = $('#demoAdmin');
     //监听导航点击
     element.on('nav', function(elem) {
         var link = $(elem).attr('wsz-href');
         $iframe.prop('src', link);
     });

 });

 $('#search').on('click', function() {
     $('.menu dd').attr('class', '');
     $('.menu li').attr('class', 'layui-nav-item');

     var data = $('.search').val();
     var $iframe = $('#demoAdmin');

     $('a').each(function() {
         if ($(this).text() == data) {
             var link = $(this).attr('wsz-href');
             $iframe.prop('src', link);
             $(this).parent('dd').addClass('layui-this');
             $(this).parents('dd').addClass('layui-nav-itemed').parents('li').addClass('layui-nav-itemed');
         }
     })
 })

 $('.layui-nav-tree').children('li').children('a').on('click', function() {
     if ($(this).parent().hasClass('layui-nav-itemed')) {
         $(this).parent().siblings('li').removeClass('layui-nav-itemed');
         $(this).parent().addClass('layui-nav-itemed');
     }
 });

 $('.layui-nav-tree a').on('click', function() {
     $('.layui-layout-left').children('li').removeClass('layui-this');
 });

 $('.layui-layout-left').on('click', function() {
     $('.layui-nav-tree').children('li').removeClass('layui-nav-itemed').children('dd').removeClass('layui-this');
 });
 //  http://www.pdmreader.com/team/demo/reportDemo.html