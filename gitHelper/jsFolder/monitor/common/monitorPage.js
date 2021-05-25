var fadeInTime=382;

//展示遮罩
function showScreenShade()
{
    $("#screenShade").css("z-index",100);
    $("#screenShade").show();
}

function beforShade(divId)
{
    showScreenShade();//

    var screenZIndex=parseInt($('#screenShade').css('z-index'));
    screenZIndex=screenZIndex+1;
    $("#" + divId).css("z-index",screenZIndex);
}


function createDivWindow(divId,height,width)
{
    beforShade(divId);

    if (width==null)
        width=height/0.618;
    if (width==0)
        width=height/0.618;

    screenWidth=$(window).width();
    screenHeight=$(window).height();
    var theWidth=width;
    var theHeight=height;
    var theLeft=(screenWidth-theWidth)/2;
    var theTop=(screenHeight-theHeight)/2;

    $( "#" + divId).css("width",theWidth);
    $( "#" + divId).css("height",theHeight);
    $( "#" + divId).css("top",theTop);
    $( "#" + divId).css("left",theLeft);

    $( "#" + divId).fadeIn(fadeInTime);
}

function closeDivWindow(divId)
{
    $("#screenShade").hide();
    $("#" + divId).hide();
}