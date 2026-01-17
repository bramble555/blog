import service from './axios'

export const getDataSum = () => {
    return service.get('/data/sum')
}

export const getUserLoginData = () => {
    return service.get('/data/login_data')
}
