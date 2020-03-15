$(document).ready(function() {
  /* 스크롤 300이하시 top 메뉴 display:none, 아닐시 display:block;*/
  $(window).on("scroll", function() {
    var winScroll = $(this).scrollTop();

    if (winScroll <= 300) {
      $("div#quickTop").css({ display: "none" });
    } else $("div#quickTop").css({ display: "block" });

    if (winScroll > 180) {
      $("nav").addClass("scroll");
    } else $("nav").removeClass("scroll");
  });

  /* top 클릭하고 상위로 이동할 때 에니메이트*/
  $("div#quickTop").on("click", function() {
    var winAnimate = $("html, body");
    $(winAnimate)
      .stop()
      .animate({ scrollTop: "0" }, 500, "linear", function() {
        //alert("haha");
      });
  });

  /* 검색창 jQuery */
  $(".btn_search").on("mouseenter", function() {
    $(".searchField").animate({ width: "220px" }, "linear");
  });
  $(".btn_search").on("mouseleave", function() {
    $(".searchField")
      .stop()
      .animate({ width: "150px" }, "linear");
  });

  /* 서브 메뉴창 */

  /*마우스를 가져갔을 때 lnb sub에 addClass를 하여 display:block;*/
  $(".pupleColor, .lnb_sub").on("mouseenter", function() {
    //console.log($(".navClose").css("display"));
    var mobileNav = $(".navClose").css("display");
    if (mobileNav == "none") {
      $(".lnb_sub").addClass("on");
    }
  });
  /*마우스를 떼었을때 lnb sub에 removeClass를 하여 display:none;*/
  $(".pupleColor, .lnb_sub").on("mouseleave", function() {
    var mobileNav = $(".navClose").css("display");
    if (mobileNav == "none") {
      $(".lnb_sub").removeClass("on");
    }
  });

  /* 서브 메뉴 이미지 변경 */

  $(".fruit a").hover(
    function() {
      $("img", this).attr("src", "./images/common/nav01_fruite_ov.png");
    },
    function() {
      $("img", this).attr("src", "./images/common/nav01_fruite.png");
    }
  );

  $(".fish a").hover(
    function() {
      $("img", this).attr("src", "./images/common/nav02_fish_ov.png");
    },
    function() {
      $("img", this).attr("src", "./images/common/nav02_fish.png");
    }
  );

  $(".meat a").hover(
    function() {
      $("img", this).attr("src", "./images/common/nav03_meat_ov.png");
    },
    function() {
      $("img", this).attr("src", "./images/common/nav03_meat.png");
    }
  );

  $(".sidedish a").hover(
    function() {
      $("img", this).attr("src", "./images/common/nav04_sidedish_ov.png");
    },
    function() {
      $("img", this).attr("src", "./images/common/nav04_sidedish.png");
    }
  );

  $(".retrofood a").hover(
    function() {
      $("img", this).attr("src", "./images/common/nav05_retro_ov.png");
    },
    function() {
      $("img", this).attr("src", "./images/common/nav05_retro.png");
    }
  );

  $(".pastaNoil a").hover(
    function() {
      $("img", this).attr("src", "./images/common/nav06_seasoning_ov.png");
    },
    function() {
      $("img", this).attr("src", "./images/common/nav06_seasoning.png");
    }
  );

  $(".beverage a").hover(
    function() {
      $("img", this).attr("src", "./images/common/nav07_bottle_ov.png");
    },
    function() {
      $("img", this).attr("src", "./images/common/nav07_bottle.png");
    }
  );

  $(".cheese a").hover(
    function() {
      $("img", this).attr("src", "./images/common/nav08_cheese_ov.png");
    },
    function() {
      $("img", this).attr("src", "./images/common/nav08_cheese.png");
    }
  );

  $(".beauty a").hover(
    function() {
      $("img", this).attr("src", "./images/common/nav09_make-up_ov.png");
    },
    function() {
      $("img", this).attr("src", "./images/common/nav09_make-up.png");
    }
  );

  $(".living a").hover(
    function() {
      $("img", this).attr("src", "./images/common/nav10_living_ov.png");
    },
    function() {
      $("img", this).attr("src", "./images/common/nav10_living.png");
    }
  );

  $(".baby a").hover(
    function() {
      $("img", this).attr("src", "./images/common/nav11_baby_ov.png");
    },
    function() {
      $("img", this).attr("src", "./images/common/nav11_baby.png");
    }
  );

  $(".brandStore a").hover(
    function() {
      $("img", this).attr("src", "./images/common/nav12_shoppingbag_ov.png");
    },
    function() {
      $("img", this).attr("src", "./images/common/nav12_shoppingbag.png");
    }
  );

  $(".MDcuration a").hover(
    function() {
      $("img", this).attr("src", "./images/common/nav13_MD_ov.png");
    },
    function() {
      $("img", this).attr("src", "./images/common/nav13_MD.png");
    }
  );

  /* 퀵메뉴 클릭시 이미지 변경*/

  $("#quickToggle img").attr({
    src: "./images/common/open-arrow.png",
    title: "퀵메뉴 열기"
  });

  $("#quickToggle").on("click", function() {
    var strImg = $("#quickToggle")
      .find("img")
      .attr("src");
    //console.log(strImg);
    //console.log(strImg.indexOf("close-arrow")); //-1
    if (strImg.indexOf("close-arrow.png") >= 0)
      $("img", this).attr("src", "./images/common/open-arrow.png");
    else $("img", this).attr("src", "./images/common/close-arrow.png");
  });

  /* 퀵 메뉴 나오기 들어가기 */
  $("#quickToggle").on("click", function() {
    $(".quick_nav").toggleClass("quick_navOn");
    $("#quickToggle").toggleClass("quick_Toggle");
  });

  /* 모바일 메뉴 추가 카테고리 */
  $(".pupleColor").click(function() {
    var winWeight = $(window).outerWidth();
    if (winWeight <= 799)
      $(".pupleColor")
        .children("ul")
        .toggleClass("mobileList");
  });

  /* 회원가입 입력사항 */
  $("#btnReset:button").click(function() {
    //alert("회원가입 양식을 다시 작성하여 주십시오. ");
    clearInput();
  });
  $("#btnSubmit:button").click(function() {
    //alert("회원가입을 축하드립니다! ");
    verification(); //호출
  });

  var verification = function() {
    if ($("#userID").val() == "") {
      alert("아이디를 입력하여 주십시오.");
      return; //처음부터 함수시
    }
    if ($("#userPwd").val() == "") {
      alert("비밀번호를 입력하세요. ");
      return;
    }
    if ($("#userrePwd").val() == "") {
      alert("비밀번호 확인을 입력하세요. ");
      return;
    }
    if ($("#userPwd").val() != $("#userrePwd").val()) {
      alert("비밀번호가 일치하지 않습니다. ");
      $("#userrePwd").focus();
      return;
    }

    alert($("#userID").val() + "님 회원가입을 축하드립니다.");
  };
  function clearInput() {
    $("#userID").val("");
    $("#userPwd").val("");
    $("#userrePwd").val("");
  }

  /* 약관 동의 */

  document.form2.allSelect.onclick = function() {
    if (this.checked) {
      document.form2.accept1.checked = true;
      document.form2.accept2.checked = true;
    } else {
      document.form2.accept1.checked = false;
      document.form2.accept2.checked = false;
    }
  };

  /* 로그인 팝업 */
  $(".login_popup").on("click", function() {
    $("#loginWrapper").addClass("on");
    $("#login_bg").addClass("on");
    $("#login_main").addClass("on");
    $("#loginClose").addClass("on");
  });

  $("#loginClose").on("click", function() {
    $("#loginWrapper").removeClass("on");
    $("#login_bg").removeClass("on");
    $("#login_main").removeClass("on");
    $("#loginClose").removeClass("on");
  });

  $(window).scroll();
});
