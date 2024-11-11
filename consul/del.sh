# 查詢所有無健康的服務，並假設 services 以逗號分隔
services="test_name.1731316064479870976,test_name.1731316073796468736,test_name.1731316234952114176,test_name.1731316235772100608,test_name.1731316243832504320"

# 遍歷每個服務並刪除
for service in $(echo "$services" | tr ',' '\n'); do
    echo "Deregistering service: $service"
    curl -X PUT http://localhost:8500/v1/agent/service/deregister/$service
done
