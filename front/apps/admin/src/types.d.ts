export { };
declare global {
  interface AccountState {
    token: string
  }
  interface RouterContext {
    readonly account: AccountState;
  }
  class Res<T = unknown> {
    readonly code: CodeType;
    readonly msg: string;
    readonly data?: T;
  }
}
