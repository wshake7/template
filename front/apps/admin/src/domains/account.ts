export interface ReqPwdLogin {
  username: string
  pwd: string
}

export interface ResLogin {
  token: string
  publicKey: string
}

export interface ReqChangePwd {
  readonly oldPwd: string
  readonly newPwd: string
}
