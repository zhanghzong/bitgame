package constants

// 私钥
var PrivateKey = []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIICWwIBAAKBgQDMUws+7NKknmImMYUsSr4DOKYVrs1s7BQzGBgkkTptjGiektUmxm3BNZq34ugF6Vob9V0vU5r0S7vfyuOTC87uFeGe+rBJf7si4kE5wsJiEBlLNZjrz0T30xHGJlf+eizYVKPkpo3012rKvHN0obBlN7iBsdiGpLGP3sPAgO2tFQIDAQABAoGAZXPyiIsUyHJwL6C1DFoMYRMWvHtwOt455WjYTAfkaBKou9wShE9Qnffc2+OJ662DdZBudZpgvV6BacyXFSNu2jvOgpZLYSZiXhd8XF9/Kpy9tmjAAf4xxrpZ1ZsqPbsgIdMsbK9ESDJXYgl5WGqYZllnz1WHboowkjeIaEuqn4ECQQDpPNUmcwNlolOa7Ubawj4nMRVcHefTrEX6xQtUP/E3ur3vVkffEjkEGBWigW07K8TrkNXR/mcXRFeJcckqJ0LhAkEA4EPZH8FDSZP6wrvliJHSC+mkFg2adg5Unrd8OxiSY8+MNoouen1BmVCsvWclUbIQ/md3iFlpJrTUhBeqSfvktQJAetooC9CZAXe3QeupXqDhzBL2hUbbTYt4cNZJWV80133tfZucz8rxbU6iVq6Fsp0jZFEtzyaJdp/w29yrcSCtAQJAVuBxwCdyFZLJ1Z5McPdsU0kTU6e5anpqtYGHEq2WKCxCuO4Wy0SyoN3rzQOkJV3Bz4vtdliMr33lxbYVNcvq8QJAT24s4oCSnDS5QAOruHiIosUmJRL5xpaiGtH0aAiYv2tKWbVl76Oz2EaOImcZ6hQIdh7XBjYE4nIa77yjkGNrig==
-----END RSA PRIVATE KEY-----
`)

// 公钥
var PublicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDMUws+7NKknmImMYUsSr4DOKYVrs1s7BQzGBgkkTptjGiektUmxm3BNZq34ugF6Vob9V0vU5r0S7vfyuOTC87uFeGe+rBJf7si4kE5wsJiEBlLNZjrz0T30xHGJlf+eizYVKPkpo3012rKvHN0obBlN7iBsdiGpLGP3sPAgO2tFQIDAQAB
-----END PUBLIC KEY-----
`)

// 系统异常错误定义
var Error = map[string]interface{}{
	"1000": map[string]interface{}{
		"res":  1000,
		"cmd": "BAD_CMD",
	},

	"1001": map[string]interface{}{
		"res":  1001,
		"cmd": "NO_CMD",
	},

	"1002": map[string]interface{}{
		"res":  1002,
		"cmd": "BAD_JWT_TOKEN",
	},

	"1003": map[string]interface{}{
		"res":  1003,
		"cmd": "ALREADY_LOGGED",
	},
}
