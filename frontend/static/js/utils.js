
function getSelectedOption(id) {
    return $(id + " option:selected").val();
}

function string2Json($textString) {

    var $jsonObj;

    if ($textString == null) {
        return null;
    }

    try {
        $jsonObj = JSON.parse($textString);
    } catch (e) {
        alert("got bad json string");
        return null;
    }

    return $jsonObj;
}



function httpPost($url, $data, $callback) {

    var http;

    try {
        http = new XMLHttpRequest();
    } catch (e) {
        try {
            http = new ActiveXObject("Msxml2.XMLHTTP");
        } catch (e) {
            try {
                http = new ActiveXObject("Microsoft.XMLHTTP");
            } catch (e) {
                alert("您的浏览器不支持AJAX！");
                return false;
            }
        }
    }

    http.onreadystatechange = function()
    {
        if (http.readyState == 4 && http.status == 200) {
            $callback(string2Json(http.responseText));
        } else if (http.readyState == 4 && http.status != 200) {
            alert("编辑API错误！")
        }
    };


    http.open("POST", $url, true);
    http.setRequestHeader("Content-Type", "application/json");
    http.send($data);
}
