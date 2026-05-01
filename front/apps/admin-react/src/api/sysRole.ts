import API from './index'

async function list() {
  await API.Get<Res>('/api/role/page', {
    cacheFor: 0,
  }).send()
}

async function create() {
  await API.Post<Res>('/api/role/create', {
    cacheFor: 0,
  }).send()
}

async function update() {
  await API.Post<Res>('/api/role/update', {
    cacheFor: 0,
  }).send()
}

async function del() {
  await API.Post<Res>('/api/role/del', {
    cacheFor: 0,
  }).send()
}

async function switchStatus() {
  await API.Post<Res>('/api/role/switch', {
    cacheFor: 0,
  }).send()
}

export const RoleApi = {
  list,
  create,
  update,
  del,
  switchStatus,
}
