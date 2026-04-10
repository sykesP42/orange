using System.Windows;
using OrangeDesktop.ViewModels;

namespace OrangeDesktop.Views;

public partial class MainWindowView : Window
{
    private readonly MainViewModel _viewModel;

    public MainWindowView(MainViewModel viewModel)
    {
        InitializeComponent();
        _viewModel = viewModel;
        DataContext = _viewModel;
        _viewModel.LogoutRequested += OnLogoutRequested;
        Loaded += async (_, _) => await _viewModel.LoadDataCommand.ExecuteAsync(null);
    }

    private void OnLogoutRequested()
    {
        Dispatcher.Invoke(() =>
        {
            DialogResult = false;
            Close();
        });
    }
}
