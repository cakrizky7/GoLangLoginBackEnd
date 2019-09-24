# GoLangLoginBackEnd
Simple Golang User Login &amp; Register using MySQL Database & JWT

=============================================================

List Route for API

<b>POST => 0.0.0.0:9000/register</b><br>
To register user<br>
Request : 
<pre style="background-color:lightblue;">
$.ajax({
    url:'http://<server-ip/domain>:9000/register',
    type:'POST',
    data:'{
        "Username":"user",
        "Password":"pass",
        "Nama_lengkap":"Your Name",
    }',
    contentType: "application/json; charset=utf-8",
    success: function(data){
        console.log(data)
    }
}) 
</pre>
Response :
<pre style="background-color:lightgreen;">
    {
        "message"   :   "done",
        "nama"      :   "Your Name,
        "username"  :   "user",
    }
</pre>
<b>POST => 0.0.0.0:9000/login</b><br>
Just send Http POST request with username & password as body data as JSON<br>
for example (AJAX):
<pre style="background-color:lightblue;">
$.ajax({
    url:'http://<server-ip/domain>:9000/login',
    type:'POST',
    data:'{
        "Username":"user",
        "Password":"pass"
    }',
    contentType: "application/json; charset=utf-8",
    success: function(data){
        console.log(data)
    }
}) 
</pre>
Response :
<pre style="background-color:lightgreen;">
    {
        "message"       :   "done",
        "nama"          :   "Your Name",
        "username"      :   "user",
        "token"         :   "<some token string>",
        "refreshToken"  :   "<some token string>",
    }
</pre>
<b>POST => 0.0.0.0:9000/relogin</b><br>
This route is for make auto login by checking token in database. if the Authorization token is not valid anymore, then you can access this url to get new JWT Token without sending username & password by sending refresh_token that you got when you attempt login. the relogin controller will check the refresh_token in database. if it is still valid, then you can get new token & refresh_token. <br>
for example (AJAX):
<pre style="background-color:lightblue;">
$.ajax({
    url:'http://<server-ip/domain>:9000/relogin',
    type:'POST',
    data:'{
        "Username":"user",
        "Token":"refresh_token"
    }',
    contentType: "application/json; charset=utf-8",
    success: function(data){
        console.log(data)
    }
}) 
</pre>
Response :
<pre style="background-color:lightgreen;">
    {
        "message"       :   "done",
        "nama"          :   "Your Name",
        "username"      :   "user",
        "token"         :   "<some token string>",
        "refreshToken"  :   "<some token string>",
    }
</pre>
<b>// Restricted<br>
GET => 0.0.0.0:9000/</b><br>
To access restricted Route add Authorization to header<br>
<pre style="background-color:lightblue;">
    'Authorization': 'Bearer ' + token
</pre>
for example (AJAX):
<pre style="background-color:lightblue;">
$.ajax({
    url:'http://<server-ip/domain>:9000/',
    type:'GET',
    headers: {
        'Authorization': 'Bearer ' + token
    },
    success: function(data){
        console.log(data)
    }
}) 
</pre>
Response :
<pre style="background-color:lightgreen;">
    {
        "data"       :   "user",
    }
</pre><br>
<b>POST => 0.0.0.0:9000/logincheck</b><br>
To manualy check JWT token validity<br>
<pre style="background-color:lightblue;">
    'Authorization': 'Bearer ' + token
</pre>
for example (AJAX):
<pre style="background-color:lightblue;">
$.ajax({
    url:'http://<server-ip/domain>:9000/logincheck',
    type:'GET',
    headers: {
        'Authorization': 'Bearer ' + token
    },
    success: function(data){
        console.log(data)
    }
}) 
</pre>
Response :
<pre style="background-color:lightgreen;">
    {
        "data"       :   "user",
    }
</pre><br>

=============================================================

To deploy

1. Git Clone https://github.com/cakrizky7/GoLangLoginBackEnd.git
2. go get 
3. Prepare the Database<br>
    3a. Add table "users"<br>
    3b. Add fields:
        <pre style="background-color:lightblue;">
        Id            varchar(128) -> Primary Key
        Username      varchar(64) -> Unique
        Password      varchar(128)
        Nama_lengkap  varchar(64)
        Token_expired text
        Created_at    timestamp
        Updated_at    timestamp
        </pre>
4. Set database connection config in "config/app.conf"
5. go run main.go

