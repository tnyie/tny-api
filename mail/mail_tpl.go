package mail

const htmlMailVerificationTemplate = `
<!DOCTYPE html>
<html>

<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width">
    <style>
        @import url('https://fonts.googleapis.com/css2?family=Roboto:wght@100;300;400;500;700;900&display=swap');

        html,
        body {
            margin: 0;
            padding: 0;
            font-size: 14px;
            font-weight: 400;
            color: 'black' !important;
            font-family: 'Roboto', sans-serif !important;
            overflow-x: hidden;
        }
    </style>
</head>

<body
    style="margin: 0;padding: 0;font-size: 14px;overflow-x: hidden;font-family: 'Roboto', sans-serif !important;">
    <div id="body"
        style="margin: 0;height: 100%%;width: 100%%;">
        <header>
        <h1>
            Verification needed for TnyIE
        </h1>
        </header>
        <p>
            Click the link below to activate your account
        </p>
        <div>
            <div">
                <a href="%s" style="color: #009688;text-decoration: none;">%s</a>
            </div>
        </div>
        <p style="margin-bottom: 2em;">
            If you did not sign up for an account on <a href="https://tny.ie" style="color: #009688">https://tny.ie</a>, then ignore this email.
        </p>
        <footer><a href="https://ui.tny.ie/tos" style="color: #009688">Terms of Service</a></footer>
    </div>
</body>
</html>
`

const textMailVerificationTemplate = `
Verification needed for TnyIE
If you did not sign up for an account on https://tny.ie, ignore this email.
Go to the following link to activate your account
%s
Terms of Service: https://ui.tny.ie/tos
`

const htmlPasswordVerificationTemplate = `
<!DOCTYPE html>
<html>

<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width">
    <style>
        @import url('https://fonts.googleapis.com/css2?family=Roboto:wght@100;300;400;500;700;900&display=swap');

        html,
        body {
            margin: 0;
            padding: 0;
            font-size: 14px;
            font-weight: 400;
            color: 'black' !important;
            font-family: 'Roboto', sans-serif !important;
            overflow-x: hidden;
        }
    </style>
</head>

<body
    style="margin: 0;padding: 0;font-size: 14px;overflow-x: hidden;font-family: 'Roboto', sans-serif !important;">
    <div id="body"
        style="margin: 0;height: 100%%;width: 100%%;">
        <header>
        <h1>
            Password Reset for TnyIE
        </h1>
        </header>
        <p>
            Click the link below to reset password
        </p>
        <div>
            <div">
                <a href="%s" style="color: #009688;text-decoration: none;">%s</a>
            </div>
        </div>
        <p style="margin-bottom: 2em;">
            If you did not request a password reset for <a href="https://tny.ie" style="color: #009688">https://tny.ie</a>, then ignore this email
        </p>
        <footer><a href="https://ui.tny.ie/tos" style="color: #009688">Terms of Service</a></footer>
    </div>
</body>
</html>
`

const textPasswordVerificationTemplate = `
Verification needed for TnyIE
If you did not sign up for an account on https://tny.ie, ignore this email.
Go to the following link to activate your account
%s
Terms of Service: https://ui.tny.ie/tos
`
