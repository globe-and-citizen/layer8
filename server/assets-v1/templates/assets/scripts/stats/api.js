export const getClientUsageStatistics = async (baseUrl, token) => {
    const resp = await window.fetch(
        baseUrl + "/api/v1/usage-stats",
        {
            method: "GET",
            headers: {
                "Content-Type": "Application/Json",
                Authorization: `Bearer ${token}`,
            },
        }
    );

    const data = await resp.json();
    return data.data
}