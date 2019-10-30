<div>
    {{.ServerName}}欢迎您
    <p>
        您的邮箱:
        <span><a href="{{.Email}}">{{.Email}}</a></span>
        验证码:
    </p>
    <p style="font-size: 32px;color: brown;">{{.Code}}</p>
    <p style="font-size: 14px;">验证码有效期为{{.Expiration}}</p>
</div>