
function get(url,data,on_success,on_failure) {

	var params = typeof data == 'string' ? data : Object.keys(data).map(
	    function(k){ return encodeURIComponent(k) + '=' + encodeURIComponent(data[k]) }
	).join('&');
	
	xhr = new XMLHttpRequest();
	xhr.onload = function() {
		if (xhr.status === 200) {
			on_success(JSON.parse(xhr.responseText));
		} else if (xhr.status !== 200) {
			on_failure(xhr.status);
		}
	};
	xhr.open('GET', url+"?"+params);
	xhr.setRequestHeader('X-Requested-With', 'XMLHttpRequest');
	xhr.send(params);
}

function post(url,data,on_success,on_failure) {

	var params = typeof data == 'string' ? data : Object.keys(data).map(
	    function(k){ return encodeURIComponent(k) + '=' + encodeURIComponent(data[k]) }
	).join('&');
	  
	xhr = new XMLHttpRequest();
	xhr.onload = function() {
		if (xhr.status === 200) {
			on_success(JSON.parse(xhr.responseText));
		} else if (xhr.status !== 200) {
			on_failure(xhr.status);
		}
	};
	xhr.open('POST', url);
	xhr.setRequestHeader('X-Requested-With', 'XMLHttpRequest');
	xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');	  
	xhr.send(params);
}

function post_json(url,data,on_success,on_failure) {

	xhr = new XMLHttpRequest();
	xhr.onload = function() {
		if (xhr.status === 200) {
			on_success(JSON.parse(xhr.responseText));
		} else if (xhr.status !== 200) {
			on_failure(xhr.status);
		}
	};
	xhr.open('POST', url);
	xhr.setRequestHeader("Content-Type", "application/json;charset=UTF-8");
	xhr.send(JSON.stringify(data));
}

