<div>
    {{.ServerName}}:
    <p>
        您的邮箱:
        <span><a href="{{.Email}}">{{.Email}}</a></span>
        验证码:
    </p>
    <p style="font-size: 32px;color: brown;">{{.Code}}</p>
    <p style="font-size: 14px;">验证码有效期为{{.Expiration}}</p>
    <p style="font-size: 14px;">接收到该邮件，说明您将要修改密码，若非本身所为，请忽略该邮件</p>
</div>