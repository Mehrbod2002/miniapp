import axios from "axios";

export const logError = async (error: any, user_id: number, username: string | null) => {
    try {
        await axios.post(`https://miniapp.dandanjan.ir/backend/log`, {
            message: error.message || "Unknown error occurred",
            stack: error.stack || "No stack trace available",
            user: { id: user_id, username: username },
        });
        console.log("Error logged successfully.");
    } catch (loggingError) {
        console.error("Failed to log error:", loggingError);
    }
};
