# 查詢所有無健康的服務，並假設 services 以逗號分隔
services="1731467813895274496,1731468000715866112,1731316234952114176"

# 遍歷每個服務並刪除
for service in $(echo "$services" | tr ',' '\n'); do
    echo "Deregistering service: $service"
    curl -X PUT http://localhost:8500/v1/agent/service/deregister/$service
done
