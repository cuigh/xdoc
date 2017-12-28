/* global $ */
/* global hljs */
function Document() {
    this.initialize = function () {
        // 处理链接
        $("#div_content a").attr("target", "_blank");
        // 代码高亮
        $('pre code').each(function (i, block) {
            hljs.highlightBlock(block);
        });
        // 绘图
        mermaid.init(undefined, ".language-mermaid");
        // 目录
        this.initTOC();
        // 顶部跳转
        this.initUp();
    }

    this.initTOC = function () {
        $('#a_toc').click(function (e) {
            $('#div_toc').toggle();
        });

        // create toc
        var $elems = $('#div_content').find("h2,h3,h4,h5,h6");
        if ($elems.length == 0) {
            $('#a_toc,#div_toc').remove();
            return;
        }

        var rootNode = { parent: null, tag: "", text: "目录", anchor: "", nodes: [], index: 1, id: "" };
        var preNode = rootNode;
        $elems.each(function (i, e) {
            var $e = $(e);
            var node = { parent: null, tag: e.tagName, text: $e.text(), anchor: "", nodes: [], index: 0, id: "" };
            if (e.tagName !== preNode.tag) {
                var n = preNode;
                while (n.tag >= node.tag) {
                    n = n.parent;
                }
                node.parent = n;
                n.nodes.push(node);
            } else {
                node.parent = preNode.parent;
                preNode.parent.nodes.push(node);
            }
            node.index = node.parent.nodes.length;
            node.id = node.parent.id === "" ? node.index.toString() : (node.parent.id + "." + node.index);
            preNode = node;

            e.id = "h." + node.id;
            $e.text(node.id + " " + node.text)
        });

        function createIndex(array, node) {
            array.push('<li><a href="#h.' + node.id + '">' + node.id + ". " + node.text + '</a>')
            if (node.nodes.length > 0) {
                array.push('<ol>');
                for (var i = 0; i < node.nodes.length; i++) {
                    createIndex(array, node.nodes[i]);
                }
                array.push('</ol>');
            }
            array.push('</li>');
        };

        var array = [];
        array.push('<ol style="padding-left:5px">');
        array.push('<li><p style="color:gray;margin-bottom:5px">目录</p></li>');
        for (var i = 0; i < rootNode.nodes.length; i++) {
            createIndex(array, rootNode.nodes[i]);
        }
        array.push('</ol>');
        $("#div_toc").append(array.join(""));
    }

    this.initUp = function() {
        var slideToTop = $("<div />");
        slideToTop.html('<i class="glyphicon glyphicon-chevron-up"></i>');
        slideToTop.css({
            position: 'fixed',
            bottom: '20px',
            right: '25px',
            width: '40px',
            height: '40px',
            color: '#eee',
            'font-size': '',
            'line-height': '40px',
            'text-align': 'center',
            'background-color': '#222d32',
            cursor: 'pointer',
            'border-radius': '5px',
            'z-index': '99999',
            opacity: '.7',
            'display': 'none'
        });
        slideToTop.on('mouseenter', function () {
            $(this).css('opacity', '1');
        });
        slideToTop.on('mouseout', function () {
            $(this).css('opacity', '.7');
        });
        $('#wrapper').append(slideToTop);
        $(slideToTop).click(function () {
            $("html,body").animate({
                scrollTop: 0
            }, 500);
        });
        $(window).scroll(function () {
            if ($(window).scrollTop() >= 150) {
                if (!$(slideToTop).is(':visible')) {
                    $(slideToTop).fadeIn(500);
                }
            }
            else {
                $(slideToTop).fadeOut(500);
            }
        });
    }
}

// 初始化
$(function () { new Document().initialize(); });
