package mail

const textTemplate = `
Verification needed for TnyIE

Go to the following link to activate your account

%s

If you did not sign up for an account on https://tny.ie, ignore this emailconst

Terms of Service: https://tny.ie/ui/tos

`

const htmlTemplate = `
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
            font-family: 'Roboto', sans-serif !important;
            color: white;
            overflow-x: hidden;
        }
    </style>
</head>

<body
    style="margin: 0 auto;padding: 0;font-size: 14px;color: white;overflow-x: hidden;font-family: 'Roboto', sans-serif !important;">
    <div id="body"
        style="margin: 0 auto;background-color: #121212;height: 100%%;position: absolute;width: 100%%;">
        <header style="width: 100%%; background-color: #272727;height: 60px;position: relative;font-size: 20px;margin-bottom: 1px solid #121212;display: block;text-align: center;line-height:60px">
            <span>Tny</span><span style="color: #009688">IE</span>
        </header>
        <h1 style="color: #fff; font-size: 18px;text-align: center;padding: 8px 0;margin: 2em 0;">
            Verification needed for TnyIE
        </h1>
        <p style="margin: auto;text-align:center;">
            Click the link below to activate your account
        </p>
        <div style="max-width: 460px; margin: 1em auto">
            <div style="text-align: center;margin: 6em 0;">
                <a href="%s" style="color: white;background-color: #009688; padding: 1.4em 2em; border-radius: 4%%;text-decoration: none;">Verify</a>
            </div>
        </div>
        <p style="text-align: center;margin-bottom: 2em;">
            If you did not sign up for an account on <a href="https://tny.ie" style="color: #009688">https://tny.ie</a>, then ignore this email.
        </p>
        <footer style="text-align: center;margin-bottom: 2em;"><a href="https://tny.ie/ui/tos" style="color: white;">Terms of Service</a></footer>
    </div>
</body>
</html>
`
