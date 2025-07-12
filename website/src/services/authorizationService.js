
import axios from 'axios'

/**
 * Check if user can connect to chat
 * @param {string} userName
 * @param {string} serverId
 * @returns {Promise<[boolean, string]>} [success, error text]
 */
export async function canConnect(userName, serverId) {
    return axios.get(`/api/chat/canConnect/${serverId}`, {
        params: { User: userName }
    })
    .then(response => processResponse(response))
    .catch(error => {
        if (error.response) {
            return processResponse(error.response)
        }
        return [false, '']
    })
}

/**
 * Process response for axios
 * @param {object} response
 * @returns {[boolean, string]} [success, error text]
 */
function processResponse(response) {

    if(response == null) {
        return [false, 'Connection to server failed']
    }

    switch (response.status)
    {
        case 200:
            return [true, '']

        case 422:
            if (response.data && response.data.error) {
                return [false, response.data.error]
            }
    }


    return [false, 'Connection to server failed']
}


// async function checkApiIsAlive() {
//     let ret
//     let url = `/api/isAlive`
//     try {
//         const response = await axios.get(url)
//         ret = response.status === 200
//     } catch (error) {
//         handleAxiosErrors(error)
//         ret = false
//     }
//     return ret
// }