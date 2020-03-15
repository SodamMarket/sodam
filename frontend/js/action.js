
$(document).ready(function(){
    

    /* bx slider */
     var main = $('.bxslider').bxSlider({ 
    mode: 'horizontal', 
    auto: true,	//자동으로 슬라이드 
    controls : false,	//좌우 화살표	
    autoControls: false,	//stop,play 
    pager:true,	//페이징 
    pause: 2500, 
    autoDelay: 0,	
    slideWidth: 1280, 
    speed: 500, 
    stopAutoOnclick:true 
         
         
}); 


    
 /* 스크롤 300이하시 top 메뉴 display:none, 아닐시 display:block;*/
  $(window).on("scroll",function(){
      
   var winScroll = $(this).scrollTop();
        if (winScroll <= 300)
         $("div#quickTop").css({"display":"none"});
        else
            $("div#quickTop").css({"display":"block"});
        
       if (winScroll>180){
           $("nav").addClass("scroll");
         
           }else
               $("nav").removeClass("scroll");
               
      
    });
    
    /* top 클릭하고 상위로 이동할 때 에니메이트*/
    $("div#quickTop").on("click",function(){
        var winAnimate = $("html, body");
        $(winAnimate).stop().animate({"scrollTop":"0"},500,'linear',function(){
            //alert("haha");
        });
        
    });


    /* 검색창 jQuery */
    $(".btn_search").on("mouseenter",function(){
        $(".searchField").animate({"width":"220px"},"linear");
    });
    $(".btn_search").on("mouseleave",function(){
        $(".searchField").stop().animate({"width":"150px"},"linear");
         
    });
    
   
    
    /*마우스를 가져갔을 때 lnb sub에 addClass를 하여 display:block;*/
    $(".pupleColor, .lnb_sub").on("mouseenter",function(){
        
        //console.log($(".navClose").css("display"));
        var mobileNav = $(".navClose").css("display");
        if(mobileNav == "none"){
            $(".lnb_sub").addClass("on");
        }
           
       });
    /*마우스를 떼었을때 lnb sub에 removeClass를 하여 display:none;*/
    $(".pupleColor, .lnb_sub").on("mouseleave",function(){
        
        var mobileNav = $(".navClose").css("display");
        if( mobileNav == "none"){
            $(".lnb_sub").removeClass("on");
        }
       });
      

/* 서브 메뉴 이미지 변경 */ 

    $(".fruit a").hover(function(){
        $("img", this).attr("src","./images/common/nav01_fruite_ov.png");},function(){
            $("img", this).attr("src","./images/common/nav01_fruite.png");
        });
     
    
    $(".fish a").hover(function(){
        $("img", this).attr("src","./images/common/nav02_fish_ov.png");},function(){
            $("img", this).attr("src","./images/common/nav02_fish.png");
        });
      

    $(".meat a").hover(function(){
        $("img", this).attr("src","./images/common/nav03_meat_ov.png");},function(){
            $("img", this).attr("src","./images/common/nav03_meat.png");
    });
      
   
    $(".sidedish a").hover(function(){
        $("img", this).attr("src","./images/common/nav04_sidedish_ov.png");},function(){
            $("img", this).attr("src","./images/common/nav04_sidedish.png");
       });
    

    $(".retrofood a").hover(function(){
        $("img", this).attr("src","./images/common/nav05_retro_ov.png");},function(){
            $("img", this).attr("src","./images/common/nav05_retro.png");
       });
    
    $(".pastaNoil a").hover(function(){
        $("img", this).attr("src","./images/common/nav06_seasoning_ov.png");},function(){
            $("img", this).attr("src","./images/common/nav06_seasoning.png");
       });

    $(".beverage a").hover(function(){
        $("img", this).attr("src","./images/common/nav07_bottle_ov.png");},function(){
            $("img", this).attr("src","./images/common/nav07_bottle.png");
       });
    
    $(".cheese a").hover(function(){
        $("img", this).attr("src","./images/common/nav08_cheese_ov.png");},function(){
            $("img", this).attr("src","./images/common/nav08_cheese.png");
       });
    

    $(".beauty a").hover(function(){
        $("img", this).attr("src","./images/common/nav09_make-up_ov.png");},function(){
             $("img", this).attr("src","./images/common/nav09_make-up.png");
       });
    
    $(".living a").hover(function(){
        $("img", this).attr("src","./images/common/nav10_living_ov.png");},function(){
            $("img", this).attr("src","./images/common/nav10_living.png");              
       });
    
    $(".baby a").hover(function(){
        $("img", this).attr("src","./images/common/nav11_baby_ov.png");},function(){
             $("img", this).attr("src","./images/common/nav11_baby.png");
       });

    $(".brandStore a").hover(function(){
        $("img", this).attr("src","./images/common/nav12_shoppingbag_ov.png");},function(){
            $("img", this).attr("src","./images/common/nav12_shoppingbag.png");
       });
    
    $(".MDcuration a").hover(function(){
        $("img", this).attr("src","./images/common/nav13_MD_ov.png");},function(){
            $("img", this).attr("src","./images/common/nav13_MD.png");
       });
    
    
    /* 퀵메뉴 클릭시 이미지 변경*/
    
    $("#quickToggle img").attr({"src":"./images/common/open-arrow.png","title":"퀵메뉴 열기"});
    
    $("#quickToggle").on("click",function(){
        var strImg = $("#quickToggle").find("img").attr("src");
        console.log(strImg);
        console.log(strImg.indexOf("close-arrow")); //-1
        if(strImg.indexOf("close-arrow.png")>=0)
            $("img",this).attr("src","./images/common/open-arrow.png");
        else
            $("img",this).attr("src","./images/common/close-arrow.png");
             
    });

    /* 퀵 메뉴 나오기 들어가기 */
    $("#quickToggle").on("click",function(){
        
        $(".quick_nav").toggleClass("quick_navOn");
        $("#quickToggle").toggleClass("quick_Toggle");
    
    });
    
      /* 모바일 메뉴 로그인 - 회원가입*/
    $(".navHamburger").click(function(){
        
        var winWeight = $(window).outerWidth();
        if (winWeight <=799){
            var mobileShow = $("#lnb #mobileShow").hasClass("on");
            console.log(mobileShow);
            if (!mobileShow)
               $("#lnb #mobileShow").addClass("on");
        }else {
           $("#lnb #mobileShow").removeClass("on");
            
        }//확인받기!!
            

    });
    
    /* 모바일 메뉴 추가 카테고리 */
    $(".pupleColor").click(function(){
        var winWeight = $(window).outerWidth();
        if (winWeight <=799)
            console.log($(".pupleColor").children("ul").toggleClass("mobileList"));
   
    }); 
 
    /* 이벤트 */
    var index= 0; 
    var max=0;
    var container;
    
    function common(){
        max=$(".event_wrapper").children().length;
        container=$(".event_wrapper");

        $(".eventArea .prevBtn").on("click",function(){
            index--;
            if(index<0)
                index=max-1
            
            animate();
        });
        $(".eventArea .nextBtn").on("click",function(){
            index++;
            if(index>max-1)
                index=0
            animate();
            
        });
    }

    function animate(){
        var moveX=index*570;
            container.animate({"margin-left":-moveX},1000);
    }
    
    $(document).ready(common);
     
    
    var index=0;
    var max2=0;
    var container2;
    
     function common2(){
        max2=$(".new_wrapper").children().length;
        container2=$(".new_wrapper");
        
        $(".newproductArea .prevBtn").on("click",function(){
            index--;
            if(index<0)
                index=max2-1
            
            animate2();
        });
        $(".newproductArea .nextBtn").on("click",function(){
            index++;
            if(index>max2-1)
                index=0
            
            animate2();
        });
         
         
         
    }
    
    function animate2(){
        var moveX=index*570;
            container2.animate({"margin-left":-moveX},1000);
    };
    
   $(document).ready(common2);

    
    
  /* 로그인 팝업 */
   $(".login_popup").on("click",function(){
       $("#loginWrapper").addClass("on");
       $("#login_bg").addClass("on");
       $("#login_main").addClass("on");
       $("#loginClose").addClass("on");

   });
    
    $("#loginClose").on("click",function(){
        $("#loginWrapper").removeClass("on");
       $("#login_bg").removeClass("on");
       $("#login_main").removeClass("on");
       $("#loginClose").removeClass("on");
    });
    
    
    var bpop = $("#bpop1").bPopup({
       modalColor:'#000',
        opacity:0.3,
        appending:true
    });
    
     $(".popupClose img").on("click",function(){
            bpop.close();
        });
    
    $(window).scroll();

});



