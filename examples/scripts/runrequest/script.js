var result = nap.run("request.yml")

if (result.Error) {
    nap.fail(result.Error.Error())
}