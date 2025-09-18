import axios from "axios";

const api = axios.create({
    baseURL: import.meta.env.VITE_API_PREFIX,
    withCredentials: false,
})


export default api