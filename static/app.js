let promptHistory = [];
t = 0;
let resp = "";
var converter = new showdown.Converter();
$(document).ready(function(){  
    var storedApiKey = localStorage.getItem("apiKey");
    if (storedApiKey) {
        $("#apiKey").val(storedApiKey);
    }
    $('#send').click(function (e) {
        e.preventDefault();
        var prompt = $("#prompt").val().trimEnd();
        var apiKey = $("#apiKey").val().trim();
        
        // Check if prompt is not empty
        if (prompt !== "") {
            $("#prompt").val("");
            autosize.update($("#prompt"));
            localStorage.setItem("apiKey", apiKey);
            promptHistory.push({ input: prompt });
            $("#printout").append(
                "<div class='prompt-message'>" +
                "<div style='white-space: pre-wrap;'>" +
                prompt +
                "</div>" +
                "<span class='message-loader js-loading spinner-border'></span>" +
                "</div>"
            );
            window.scrollTo({ top: document.body.scrollHeight, behavior: 'smooth' });
            run(prompt, apiKey);
            $(".js-logo").addClass("active");
        } else {
            // Handle case where prompt is empty
            console.error("Empty prompt. Please enter a valid prompt.");
            // Optionally, display a message to the user
        }
    });     
    $('#prompt').keypress(function(event){        
        var keycode = (event.keyCode ? event.keyCode : event.which);
        if((keycode == 10 || keycode == 13) && event.ctrlKey){
            $('#send').click();
            return false;
        }
    });       
    autosize($('#prompt'));    
});  

function run(prompt, apiKey, action = "/run") {
    let t = 0;

    function myTimer() {
        t++;
    }

    const myInterval = setInterval(myTimer, 1000);

    $.ajax({
        url: action,
        method: "POST",
        data: JSON.stringify({ input: prompt, apiKey: apiKey }),
        contentType: "application/json; charset=utf-8",
        dataType: "json",
        success: function (data) {
            console.log("Successfully fetched history:", data);
            promptHistory = data.history;

            if (data.error) {
                console.error("Error:", data.error);
                // Handle error if needed
            } else {
                    var jsonData = data;
                    // Rest of your code
        
                $("#printout").append(
                    "<div class='px-3 py-3'>" +
                    "<div style='white-space: pre-wrap;'>" +
                    converter.makeHtml(data.response) +
                    "</div>" +
                    " <small class='timer'>(" + t + "s)</small> " +
                    "</div>"
                );
            }

            fetchAndDisplayHistory(apiKey);
        },
        error: function (jqXHR, textStatus, errorThrown) {
            console.error("Error:", jqXHR.status, textStatus, errorThrown);
            $("#printout").append(
                "<div class='text-danger response-message'>" +
                "<div style='white-space: pre-wrap;'>" +
                "There is a problem answering your question. Please check the command line output." +
                "</div>" +
                " <small class='timer'>(" + t + "s)</small> " +
                "</div>"
            );
        },
        complete: function () {
            clearInterval(myInterval);
            t = 0;
            $(".js-loading").removeClass("spinner-border");
            window.scrollTo({ top: document.body.scrollHeight, behavior: 'smooth' });
            hljs.highlightAll();
        }
    });
}

function fetchAndDisplayHistory(apiKey) {
    $.ajax({
        url: "/fetchHistory", // Use the new endpoint for fetching history
        method: "POST",
        contentType: "application/json",
        data: JSON.stringify({
            apiKey: apiKey,
        }),
        success: function (response) {
            console.log("Received chat history response:", response);

            // Update the promptHistory variable
            promptHistory = response.history;

            // Display chat history
            for (let i = 0; i < promptHistory.length; i++) {
                const chatEntry = promptHistory[i];
                $("#printout").append(
                    "<div class='px-3 py-3'>" +
                    "<div style='white-space: pre-wrap;'>" +
                    "<strong>User:</strong> " + chatEntry.input + "<br>" +
                    "<strong>Bot:</strong> " + chatEntry.response +
                    "</div>" +
                    "</div>"
                );
            }

            // Scroll to the bottom after displaying history
            window.scrollTo({ top: document.body.scrollHeight, behavior: 'smooth' });

            // Add class to logo (if needed)
            $(".js-logo").addClass("active");
        },
        error: function (error) {
            console.error("Error fetching chat history:", error);
        },
    });
}

