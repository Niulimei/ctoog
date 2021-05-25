function zoomImg(obj)
{
  var url=obj.src; 
  var processLayerId=layer.open({
            type: 1
            ,title: false //不显示标题栏
            ,closeBtn: true
            ,area: '92%'
            ,shade: 0.8
            ,id: 'LAY_layuipro' //设定一个id，防止重复弹出
            ,btnAlign: 'c'
            ,moveType: 1 //拖拽模式，0或者1
            ,content: '<img src="' + url + '" style="width:100%" >'
            ,success: function(layero){
              ;
            }
          });
}       