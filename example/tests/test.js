const assert = require('assert')
const axios = require('axios')

const MAIN_ROUTE = 'http://localhost:8080'
const SIGNUP_ROUTE = '/signup'

var randomName = size => {
	var text = ""
	var possible = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	for( var i=0; i < size; i++ )
		text += possible.charAt(Math.floor(Math.random() * possible.length))

	return text + '@gmail.com'
}

var randomPassword = size => {
    var text = ""
    var possible = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

    for( var i=0; i < size-2; i++ )
        text += possible.charAt(Math.floor(Math.random() * possible.length))

    return text+'10'
}

describe("Test signup", () => {
    it("Creates new user", done => {
    	axios.request({
    		url: MAIN_ROUTE + SIGNUP_ROUTE,
    		method: 'post',
    		data: {
    			name: randomName(7),
    			password: randomPassword(10)
    		}
    	}).then(response => {
    		var resp = response.data;
    		if (resp.Error) {
    			done(new Error(resp.Error))
    			return
    		}
    		done()
    	}).catch(err => done(err))
    })
})