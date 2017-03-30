const assert = require('assert')
const axios = require('axios')

const MAIN_ROUTE = 'http://localhost:8080'
const SIGNUP_ROUTE = '/signup'
const LOGIN_ROUTE = '/login'
const PRIVATE_ROUTE = '/priv'

const charSequence = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

var randomName = size => {
	var text = ""
	for( var i=0; i < size; i++ )
		text += charSequence.charAt(Math.floor(Math.random() * charSequence.length))

	return text + '@gmail.com'
}

var randomPassword = size => {
    var text = ""
    for( var i=0; i < size-2; i++ )
        text += charSequence.charAt(Math.floor(Math.random() * charSequence.length))

    return text+'10'
}

describe("Test signup/login", () => {
    var name = randomName(7)
    var password = randomPassword(10)
    var token
    it("Creates new user", done => {
    	axios.request({
    		url: MAIN_ROUTE + SIGNUP_ROUTE,
    		method: 'post',
    		data: {
    			name,
    			password
    		}
    	}).then(response => {
    		var resp = response.data;
    		if (resp.Error) {
    			done(new Error(resp.Error))
    			return
    		}
            console.log('User created:', name)
    		done()
    	}).catch(err => done(err))
    })

    it("Retrieves token for a given credentials", done => {
        axios.request({
            url: MAIN_ROUTE + LOGIN_ROUTE,
            method: 'post',
            data: {
                name,
                password
            }
        }).then(response => {
            var resp = response.data;
            if (resp.Error) {
                console.log('User:', name)
                console.log('Password:', password)
                done(new Error(resp.Error))
                return
            }
            token = resp.Data.token
            console.log('Username:', resp.Data.name)
            console.log('Token:', token)
            done()
        }).catch(err => done(err))
    })

    it("Tries to access private route using valid token", done => {
        axios.request({
            url: MAIN_ROUTE + PRIVATE_ROUTE,
            method: 'get',
            headers: {
                Authorization: token
            }
        }).then(response => {
            console.log(response.data)
            done()
        }).catch(err => done(err))
    })

    it("Tries to access private route using invalid token", done => {
        var invalidToken = token
        var maxMutations = 3;
        for (let i = 0; i < maxMutations; i++) {
            let randIndex = Math.floor(Math.random() * invalidToken.length) + 1
            let newChar = charSequence.charAt(Math.floor(Math.random() * charSequence.length))
            invalidToken = 
                invalidToken.slice(0, randIndex-1)
                .concat(newChar)
                .concat(invalidToken.slice(randIndex+1, invalidToken.length))    
        }

        console.log('Invalid token:', invalidToken)
        axios.request({
            url: MAIN_ROUTE + PRIVATE_ROUTE,
            method: 'get',
            headers: {
                Authorization: invalidToken
            }
        }).then(response => {
            if (response.status == 401) {
                done()
                return
            }
            done(new Error('Unauthorized access!'))
        }).catch(err => {
            if (err) done()
        })
    })
})