jQuery(document).ready(function() {

    $("#colorpick1").spectrum({
        showInitial: true,
        showInput: true,
        showButtons: true,
        preferredFormat: "hex6",
        change: function(color) {
            $("#colorpick1").val(color.toHex());
        },

        move: function(color) {
            $("#colorpick1").val(color.toHex());
        }
    });

    $("#colorpick2").spectrum({
        showInitial: true,
        showInput: true,
        showButtons: true,
        preferredFormat: "hex6",
        change: function(color) {
            $("#colorpick2").val(color.toHex());
        },

        move: function(color) {
            console.log("move:", color.toHex());
            $("#colorpick2").val(color.toHex());
        }
    });
});
