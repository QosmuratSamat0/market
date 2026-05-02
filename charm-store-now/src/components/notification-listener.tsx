import { useEffect } from "react";
import { toast } from "sonner";
import { useNotifications } from "@/hooks/use-notifications";

export function NotificationListener() {
  const { addNotification } = useNotifications();

  useEffect(() => {
    // Determine the SSE URL
    // Since we are using an API Gateway, we can use the relative path /notifications/events
    const sseUrl = "/notifications/events";

    console.log("Connecting to notifications at:", sseUrl);
    const eventSource = new EventSource(sseUrl);

    eventSource.onmessage = (event) => {
      console.log("Notification received:", event.data);
      try {
        const data = JSON.parse(event.data);
        if (data.status === "success" && data.order_id) {
          const title = "Payment Successful";
          const message = `Order ID: ${data.order_id.substring(0, 8)}... has been processed.`;
          
          addNotification({
            title,
            message,
            type: "success",
          });

          toast.success(message, {
            description: "Your order has been processed and is being prepared.",
          });
        } else {
          addNotification({
            title: "New Notification",
            message: data.message || "New update received",
            type: "info",
          });
          toast.info(data.message || "New update received");
        }
      } catch (e) {
        addNotification({
          title: "Update",
          message: event.data || "New notification received!",
          type: "info",
        });
        toast.info(event.data || "New notification received!");
      }
    };

    eventSource.onopen = () => {
      console.log("SSE connection established");
    };

    eventSource.onerror = (error) => {
      console.error("SSE connection error:", error);
    };

    return () => {
      console.log("Closing SSE connection");
      eventSource.close();
    };
  }, []);

  return null;
}
