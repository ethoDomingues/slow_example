
var schema = "http://"
var port =":5443"
if (location.protocol == "https:") {
    schema = "https://"
    port = ""
}

const HOST = "boatchazul.com.br"
const HOSTCDN = `${schema}cdn.${HOST}`
const HOSTAPI = `${schema}api.${HOST}`
const HOSTAUTH = `${schema}auth.${HOST}${port}`