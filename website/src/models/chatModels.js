
/** User of chat */
export class User {
    /** @type {string} */Id
    /** @type {string} */Name
    /** @type {UserRole} */Role
    /** @type {string} */DateTime
    /** @type {string} */Color

    constructor(data) {
        this.Id = data.Id || ''
        this.Name = data.Name || ''
        this.Role = new UserRole(data.Role || {})
        this.DateTime = data.DateTime || ''
        this.Color = data.Color || ''
    }

    /**
     * @param {string}  permissionName
     * @return {boolean}
     * * */
    havePermission(permissionName){
        // All access for admin
        if (this.Role.IsAdmin){
            return true
        }
        return Array.isArray(this.Role.Permissions) &&
            this.Role.Permissions.includes(permissionName)
    }
}

/** Role of user */
export class UserRole {
    /** @type {string} */Name
    /** @type {boolean} */IsAdmin
    /** @type {string[]} */Permissions

    constructor(data) {
        this.Name = data.Name || ''
        this.IsAdmin = data.IsAdmin || false
        this.Permissions = Array.isArray(data.Permissions) ? data.Permissions : []
    }
}

/** Message from Server */
export class Message {
    /** @type {string} */Type
    /** @type {string} */User
    /** @type {string|User[]} */MessageData
    /** @type {string} */DateTime

    constructor(data) {
        this.Type = data.type || ''
        this.User = data.user || ''
        this.MessageData = data.message || ''
        this.DateTime = data.dateTime || ''
    }
}

/** Message type from server */
export const MessageType = Object.freeze({
    message: 'Message',
    userList: 'UsersList',
    userNameChanged: 'UserNameChanged',
})

/** Permission type for users*/
export const PermissionType = Object.freeze({
    sendMessage: 'SendMessage',
})