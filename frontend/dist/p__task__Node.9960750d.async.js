(window.webpackJsonp=window.webpackJsonp||[]).push([[8],{B5KB:function(m,o,e){},swY4:function(m,o,e){"use strict";e.r(o);var U=e("g9YV"),O=e("wCAj"),a=e("k1fw"),g=e("9og8"),h=e("tJVT"),P=e("WmNS"),d=e.n(P),s=e("q1tI"),C=e.n(s),D=e("B5KB"),I=e.n(D),j=e("CLrh"),n=e("nKUr"),L=e.n(n),M=[{title:"\u8282\u70B9\u7F16\u53F7",dataIndex:"id",key:"id"},{title:"\u8282\u70B9IP",dataIndex:"workerUrl",key:"workerUrl"},{title:"\u5F53\u524D\u72B6\u6001",dataIndex:"status",key:"status"},{title:"\u5F53\u524D\u4EFB\u52A1\u6570",dataIndex:"taskCount",key:"taskCount"}],f=function(){var p=Object(s.useState)({pageSize:10,pageNum:1,total:8,workerList:[]}),c=Object(h.a)(p,2),t=c[0],E=c[1];Object(s.useEffect)(function(){(function(){var r=Object(g.a)(d.a.mark(function l(){var i;return d.a.wrap(function(_){for(;;)switch(_.prev=_.next){case 0:return _.next=2,j.b.getWorkList(t.total,0);case 2:i=_.sent,E(Object(a.a)(Object(a.a)({},t),{},{workerList:i.workerInfo,total:i.count}));case 4:case"end":return _.stop()}},l)}));function u(){return r.apply(this,arguments)}return u})()()},[]);var v=function(u,l){E(Object(a.a)(Object(a.a)({},t),{},{pageNum:u,pageSize:l||t.pageSize}))},b={showSizeChanger:!0,showQuickJumper:!1,showTotal:function(){return"\u5171".concat(t.total,"\u6761")},pageSize:t.pageSize,current:t.pageNum,total:t.total,onChange:v};return Object(n.jsxs)("div",{children:[Object(n.jsx)("p",{className:"nodeTitle",children:"\u4EFB\u52A1\u6267\u884C\u8282\u70B9\u5217\u8868"}),Object(n.jsx)(O.a,{columns:M,dataSource:t.workerList,pagination:b})]})};o.default=f}}]);