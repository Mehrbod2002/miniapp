import React, { useEffect, useState } from "react";
import WebApp from "@twa-dev/sdk";
import axios from "axios";

export default function Home() {
  const [telegramData, setTelegramData] = useState<any>(null);
  const [points, setPoints] = useState(0);

  useEffect(() => {
    WebApp.expand();
    WebApp.MainButton.setText("Add Points");
    WebApp.MainButton.show();

    WebApp.MainButton.onClick(() => {
      addPoints();
    });

    const userData = WebApp.initDataUnsafe.user;
    setTelegramData(userData);

    if (userData) {
      fetchPoints(userData.id);
    }
  }, []);

  const fetchPoints = async (telegramId: number) => {
    try {
      const response = await axios.get("http://localhost:8080/points", {
        params: { telegram_id: telegramId },
      });
      setPoints(response.data.points);
    } catch (error) {
      console.error("Error fetching points:", error);
    }
  };

  const addPoints = async () => {
    if (!telegramData) return;
    try {
      await axios.post("http://localhost:8080/points", {
        telegram_id: telegramData.id,
      });
      fetchPoints(telegramData.id);
    } catch (error) {
      console.error("Error adding points:", error);
    }
  };

  return (
    <div style={{ textAlign: "center", marginTop: "50px" }}>
      <h1>Telegram Mini App</h1>
      {telegramData ? (
        <>
          <p>
            Welcome, <strong>{telegramData.first_name}</strong>!
          </p>
          <p>Your Points: {points}</p>
        </>
      ) : (
        <p>Loading Telegram Data...</p>
      )}
    </div>
  );
}
